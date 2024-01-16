import { useEffect, useState } from "react";

/**
 * Combination of 'useState' and 'useEffect' for the common pattern of wanting
 * to have access to something produced by an async function (i.e: returning
 * a promise) inside a React component (which cannot itself be async).
 */
export const useAsync = <T,>(f: () => Promise<T | null>): T | null => {
  const [v, setV] = useState<T | null>(null);
  const asyncSetV = async () => setV(await f());

  useEffect(() => {
    asyncSetV();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return v;
};
