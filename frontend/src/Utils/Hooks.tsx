import { act } from "@testing-library/react";
import { EffectCallback, useEffect, useState } from "react";
import { discard } from "./Utils";

/**
 * Daniel named this
 */
export const useJarJar = <T,>(
  f: () => Promise<T | null>
): [T | null, EffectCallback] => {
  const [v, setV] = useState<T | null>(null);
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
export const useAsync = <T,>(p: Promise<T | null>): T | null =>
  useJarJar(() => p)[0];

/**
 * Fire the callback when the component first renders, and never again!
 */
export const useOnce = (f: EffectCallback) => {
  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(f, []);
};
