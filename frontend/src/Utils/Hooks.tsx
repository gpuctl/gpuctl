import { act } from "@testing-library/react";
import { EffectCallback, useEffect, useState } from "react";
import { Validated, Validation, discard, loading } from "./Utils";

/**
 * Daniel named this
 */
export const useJarJar = <T,>(
  f: () => Promise<Validated<T>>
): [Validation<T>, EffectCallback] => {
  const [v, setV] = useState<Validation<T>>(loading());
  const updateV = discard(async () => {
    const x = await f();
    act(() => setV(x));
  });

  useOnce(updateV);

  return [v, updateV];
};

/**
 * Combination of 'useState' and 'useEffect' for the common pattern of wanting
 * to await a Promise inside a React component.
 *
 * Returns 'null' until the promise returns.
 */
export const useAsync = <T,>(p: Promise<Validated<T>>): Validation<T> =>
  useJarJar(() => p)[0];

/**
 * Fire the callback when the component first renders, and never again!
 */
export const useOnce = (f: EffectCallback) => {
  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(f, []);
};
