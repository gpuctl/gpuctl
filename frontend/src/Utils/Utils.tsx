/**
 * Fires an asynchronous function but doesn't wait for the result
 */
export const discard = <T,>(f: () => Promise<T>) => {
  return () => {
    f();
  };
};
