create or replace function unapprove_player(p_game_id bigint, p_user_id uuid)
returns void as $$
begin
    -- Update approval status in game_approvals to UNAPPROVED
    update game_approvals
    set status = 'UNAPPROVED'
    where game_id = p_game_id and user_id = p_user_id;
    
    -- Ensure the update affected a row
    if not found then
        raise exception 'Approval record not found for game_id % and user_id %', p_game_id, p_user_id;
    end if;
    
    -- Delete the player from the game, ensuring they exist first
    delete from game_players
    where game_id = p_game_id and user_id = p_user_id;
    
    if not found then
        raise exception 'Player record not found for game_id % and user_id %', p_game_id, p_user_id;
    end if;
end;
$$ language plpgsql;
