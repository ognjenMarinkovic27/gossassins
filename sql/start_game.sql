create or replace function start_game(p_game_id bigint)
returns void as $$
declare
    player_ids uuid[];
    shuffled_ids uuid[];
    i int;
begin
    -- Fetch all player user_ids for the given game_id
    select array_agg(user_id) into player_ids
    from game_players
    where game_id = p_game_id;
    
    if array_length(player_ids, 1) is null or array_length(player_ids, 1) < 2 then
        raise exception 'Not enough players to form a cycle';
    end if;
    
    -- Shuffle the array to randomize target assignment
    select array_agg(user_id order by random()) into shuffled_ids
    from unnest(player_ids) as user_id;
    
    -- Assign kill codes and targets in a cycle
    for i in 1..array_length(shuffled_ids, 1) loop
        update game_players
        set kill_code = upper(substr(md5(random()::text), 1, 7)),
            target_id = shuffled_ids[(i % array_length(shuffled_ids, 1)) + 1]
        where game_id = p_game_id and user_id = shuffled_ids[i];
    end loop;

    -- Update game status to RUNNING
    update games
    set state = 'RUNNING'
    where id = p_game_id;
end;
$$ language plpgsql;