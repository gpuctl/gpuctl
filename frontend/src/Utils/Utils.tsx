import {GPUStats} from "../App"

/** Create an array of numbers that span between a given minimum and maximum */
export const range = (min: number, max: number) =>
  Array.from(Array(max - min).keys()).map((x) => x + min);

/**
 * Initialise an array of a given size, filled with elements using f (which is
 * given access to the index of the element is creating)
 */
export const makeArr = <T,>(size: number, f: (i: number) => T) =>
  range(0, size).map(f);

/**
 * Fires an asynchronous function but doesn't wait for the result
 */
export const discard = <T,>(f: () => Promise<T>) => {
  return () => {
    f();
  };
};

export const inlineLog = <T,>(x: T): T => {
  console.log(x);
  return x;
};

export const mapNullable = <T, U>(x: T | null, f: (x: T) => U): U | null =>
  x == null ? null : f(x);

enum VTag {
  Success = "Success",
  Loading = "Loading",
  Failure = "Failure",
}

// Note that a:
// {
//   data: T;
//   loading: false;
//   error: null;
// } | {
//   data: null;
//   loading: false;
//   error: Error;
// }} | {
//   data: null;
//   loading: true;
//   error: null;
// }-style design might seem more natural, but TypeScript's flow typing does
// not appear to be up to the challenge of reasoning from data == null towards
// loading: false and error: null

export type Success<T> = {
  tag: VTag.Success;
  data: T;
};

export type Failure = {
  tag: VTag.Failure;
  error: Error;
};

export type Loading = {
  tag: VTag.Loading;
};

export type Validated<T> = Success<T> | Failure;

export type Validation<T> = Validated<T> | Loading;

export const success = <T,>(x: T): Success<T> => ({
  tag: VTag.Success,
  data: x,
});

export const failure = (e: Error): Failure => ({
  tag: VTag.Failure,
  error: e,
});

export const loading = (): Loading => ({
  tag: VTag.Loading,
});

type ValidationMotive<T, U> = {
  success: (x: T) => U;
  loading: () => U;
  failure: (e: Error) => U;
};

/**
 * Eliminate a validation
 */
export function validationElim<T, U>(
  v: Validation<T>,
  motive: ValidationMotive<T, U>
): U {
  switch (v.tag) {
    case VTag.Success: {
      return motive.success(v.data);
    }
    case VTag.Failure: {
      return motive.failure(v.error);
    }
    case VTag.Loading: {
      return motive.loading();
    }
  }
}

export function isFree(s : GPUStats): Boolean {
  return s.gpu_util < 5
}

export const mapSuccess = <T, U>(
  v: Validation<T>,
  f: (x: T) => U
): Validation<U> => (v.tag === VTag.Success ? success(f(v.data)) : v);

/**
 * This function is only really convenient in postfix position IMO, but .methods
 * require opting into using interfaces and objects which I all find quite ugly
 */
export const orElse = <T, U>(v: Validation<T>, e: () => U) =>
  v.tag === VTag.Success ? v.data : e();
