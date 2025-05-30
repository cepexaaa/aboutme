map_build(X, X).

map_get([(Key, Value)|_], Key, Value).
map_get([_|Tail], Key, Value) :-
    map_get(Tail, Key, Value).

map_replace([], _, _, []).
map_replace([(K, _)|Tail], K, V, [(K, V)|Result]) :-
    map_replace(Tail, K, V, Result).
map_replace([(K1, V1)|Tail], K, V, [(K1, V1)|Result]) :-
    K \= K1,
    map_replace(Tail, K, V, Result).

map_lastKey([(Key, _)], Key).
map_lastKey([_|Tail], Key) :-
    map_lastKey(Tail, Key).

map_lastValue([(_, Value)], Value).
map_lastValue([_|Tail], Value) :-
    map_lastValue(Tail, Value).

map_lastEntry([Entry], Entry).
map_lastEntry([_|Tail], Entry) :-
    map_lastEntry(Tail, Entry).
