import { act } from "@testing-library/react";
import { useEffect, useState } from "react";

/**
 * Combination of 'useState' and 'useEffect' for the common pattern of wanting
 * to await a Promise inside a React component.
 *
 * Returns 'null' until the promise returns.
 */
export const useAsync = <T,>(p: Promise<T | null>): T | null => {
  const [v, setV] = useState<T | null>(null);
  const asyncSetV = async () => {
    const x = await p;
    act(() => setV(x));
  };

  useEffect(() => {
    asyncSetV();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return v;
};
