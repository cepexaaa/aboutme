composite(N) :- \+ prime(N).
prime(2).
prime(N) :-
    N > 1, 
    N mod 2 =\= 0,
    prime_helper(N, 3).
prime_helper(N, I) :-
    I * I > N.
prime_helper(N, I) :-
    N mod I =\= 0, 
    I2 is I + 2,
    prime_helper(N, I2).

prime_divisors(1, _, []):- !.
prime_divisors(N, Divisors) :-
    prime_divisors(N, 2, Divisors).

prime_divisors(N, P, [P|Divisors]) :-
    P =< N,
    N mod P =:= 0,
    N1 is N // P,
    prime_divisors(N1, P, Divisors), !.
prime_divisors(N, P, Divisors) :-
		P =< N,
    NextP is P + 1,
    prime_divisors(N, NextP, Divisors), !.


prime_index(P, N) :-
		prime(P),
    prime_index_helper(2, P, 1, N), !.

prime_index_helper(P, P, Count, Count) :- !.
prime_index_helper(I, P, Count, N) :-
    (prime(I) -> Count1 is Count + 1 ; Count1 = Count),
    I1 is I + 1,
    prime_index_helper(I1, P, Count1, N), !.


