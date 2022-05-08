import { useEffect, useRef } from "react";

export function useTimeout(callback, delay) {
  const savedCallback = useRef() as any;

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  useEffect(() => {
    function tick() {
      savedCallback.current();
    }

    if (delay !== null) {
      const id = setTimeout(tick, delay);
      return () => {
        clearInterval(id);
      };
    }
  }, [callback, delay]);
}
