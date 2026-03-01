{-# LANGUAGE InstanceSigs #-}

module HW3.T4
  ( State(..)
  , mapState
  , wrapState
  , joinState
  , modifyState
  , Prim(..)
  , Expr(..)
  , eval
  ) where

import           HW3.T1 (Annotated (..))

data State s a = S { runS :: s -> Annotated s a }

mapState :: (a -> b) -> State s a -> State s b
mapState f (S g) = S (\s -> case g s of
    x :# s' -> f x :# s')

wrapState :: a -> State s a
wrapState a = S (\s -> a :# s)

joinState :: State s (State s a) -> State s a
joinState (S f) = S (\s -> case f s of
    (S g) :# s' -> g s')

modifyState :: (s -> s) -> State s ()
modifyState f = S (\s -> () :# f s)

instance Functor (State s) where
    fmap :: (a -> b) -> State s a -> State s b
    fmap = mapState

instance Applicative (State s) where
    pure :: a -> State s a
    pure = wrapState

    (<*>) :: State s (a -> b) -> State s a -> State s b
    sf <*> sx = do
        f <- sf
        x <- sx
        return (f x)

instance Monad (State s) where
    (>>=) :: State s a -> (a -> State s b) -> State s b
    ma >>= f = joinState (mapState f ma)


data Prim a =
    Add a a      -- (+)
  | Sub a a      -- (-)
  | Mul a a      -- (*)
  | Div a a      -- (/)
  | Abs a        -- abs
  | Sgn a        -- signum
  deriving (Show)

data Expr =
    Val Double
  | Op (Prim Expr)
  deriving (Show)

instance Num Expr where
    x + y = Op (Add x y)
    x * y = Op (Mul x y)
    x - y = Op (Sub x y)
    abs x = Op (Abs x)
    signum x = Op (Sgn x)
    fromInteger x = Val (fromInteger x)
    negate x = Op (Sub (Val 0) x)

instance Fractional Expr where
    x / y = Op (Div x y)
    fromRational x = Val (fromRational x)

eval :: Expr -> State [Prim Double] Double
eval (Val x) = pure x

eval (Op (Add x y)) = do
    a <- eval x
    b <- eval y
    let result = a + b
    modifyState (Add a b :)
    pure result

eval (Op (Sub x y)) = do
    a <- eval x
    b <- eval y
    let result = a - b
    modifyState (Sub a b :)
    pure result

eval (Op (Mul x y)) = do
    a <- eval x
    b <- eval y
    let result = a * b
    modifyState (Mul a b :)
    pure result

eval (Op (Div x y)) = do
    a <- eval x
    b <- eval y
    let result = a / b
    modifyState (Div a b :)
    pure result

eval (Op (Abs x)) = do
    a <- eval x
    let result = abs a
    modifyState (Abs a :)
    pure result

eval (Op (Sgn x)) = do
    a <- eval x
    let result = signum a
    modifyState (Sgn a :)
    pure result

