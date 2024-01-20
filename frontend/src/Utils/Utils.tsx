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
  error: null;
};

export type Failure = {
  tag: VTag.Failure;
  data: null;
  error: Error;
};

export type Loading = {
  tag: VTag.Loading;
  data: null;
  error: null;
};

export type Validated<T> = Success<T> | Failure;

export type Validation<T> = Validated<T> | Loading;

export const success = <T,>(x: T): Success<T> => ({
  tag: VTag.Success,
  data: x,
  error: null,
});

export const failure = (e: Error): Failure => ({
  tag: VTag.Failure,
  data: null,
  error: e,
});

export const loading = (): Loading => ({
  tag: VTag.Loading,
  data: null,
  error: null,
});

/**
 * Eliminate a validation
 */
export function velim<T, U>(
  v: Validation<T>,
  success: (x: T) => U,
  otherwise: () => U,
  error: ((e: Error) => U) | null = null
) {
  switch (v.tag) {
    case VTag.Success: {
      return success(v.data);
    }
    case VTag.Failure: {
      return error === null ? otherwise() : error(v.error);
    }
    case VTag.Loading: {
      return otherwise();
    }
  }
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
