import { DependencyList, useEffect } from "react"

export const useDisposableEffect = (cb: (stack: DisposableStack) => void, deps?: DependencyList) => {
    useEffect(() => {
        const stack = new DisposableStack();
        try {
            cb(stack);
        } catch (e) {
            stack.dispose();
            throw e
        }

        return () => {
            stack.dispose();
        }
    }, deps)
}