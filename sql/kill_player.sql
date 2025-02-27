create or replace function kill_player(p_game_id bigint, p_killer_id uuid, p_kill_code text)
returns void as $$
declare
    v_target_id uuid;
    v_kill_code text;
    v_new_target_id uuid;
begin
    -- Find the killer's target and check the kill code
    select target_id into v_target_id
    from game_players
    where game_id = p_game_id and user_id = p_killer_id
    for update;
    
    if v_target_id is null then
        raise exception 'NO_TARGET';
    end if;

    -- Get the target's target and kill code
    select target_id, kill_code into v_new_target_id, v_kill_code
    from game_players
    where game_id = p_game_id and user_id = v_target_id;
    
    -- Verify the kill code matches
    if p_kill_code != v_kill_code then
        raise exception 'INVALID_CODE';
    end if;
    
    -- Mark the target as DEAD
    update game_players
    set status = 'DEAD', target_id = null
    where game_id = p_game_id and user_id = v_target_id;

    if v_new_target_id = p_killer_id then
        update games
        set state = 'DONE'
        where id = p_game_id;
    end if;
    
    -- Assign the new target to the killer, only if the target had a next target
    if v_new_target_id is not null then
        update game_players
        set target_id = v_new_target_id
        where game_id = p_game_id and user_id = p_killer_id;
    end if;
end;
$$ language plpgsql;
