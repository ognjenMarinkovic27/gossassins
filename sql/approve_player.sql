create or replace function approve_player(p_game_id bigint, p_user_id uuid)
returns void as $$
begin
    -- Update approval status in join_requests to APPROVED
    update join_requests
    set status = 'APPROVED'
    where game_id = p_game_id and user_id = p_user_id;
    
    -- Ensure the update affected a row
    if not found then
        raise exception 'Approval record not found for game_id % and user_id %', p_game_id, p_user_id;
    end if;
    
    -- Insert the player into the game with default values
    insert into game_players (game_id, user_id, kill_code, target_id, status)
    values (p_game_id, p_user_id, NULL, NULL, 'ALIVE')
    on conflict (game_id, user_id) do update
    set status = 'ALIVE', target_id = NULL;
end;
$$ language plpgsql;