create or replace function kill_player(p_game_id bigint, p_killer_id uuid, p_kill_code text)
returns void as $$
declare
    v_target_id uuid;
    v_new_target_id uuid;
begin
    -- Find the killer's target and check the kill code
    select target_id into v_target_id
    from game_players
    where game_id = p_game_id and user_id = p_killer_id
    for update;
    
    if v_target_id is null then
        raise exception 'Killer has no target';
    end if;
    
    -- Verify the kill code matches
    if not exists (
        select 1 from game_players
        where game_id = p_game_id and user_id = v_target_id and kill_code = p_kill_code
        for update
    ) then
        raise exception 'Invalid kill code';
    end if;
    
    -- Get the target's target
    select target_id into v_new_target_id
    from game_players
    where game_id = p_game_id and user_id = v_target_id;
    
    -- Mark the target as DEAD
    update game_players
    set status = 'DEAD', target_id = null
    where game_id = p_game_id and user_id = v_target_id;
    
    -- Assign the new target to the killer, only if the target had a next target
    if v_new_target_id is not null then
        update game_players
        set target_id = v_new_target_id
        where game_id = p_game_id and user_id = p_killer_id;
    end if;
end;
$$ language plpgsql;
