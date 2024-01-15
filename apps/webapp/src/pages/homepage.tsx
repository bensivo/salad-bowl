import { useAppDispatch, useAppSelector } from "../store/hooks"
import { counterActions } from "../store/slices/counter";

export function Homepage() {

    const count = useAppSelector(s => s.counter.count);
    const dispatch = useAppDispatch();

    return (
        <>
            <h1>Vite + React</h1>
            <div className="card">
                <button onClick={() => {
                    dispatch(counterActions.incrementBy(2));
                }}>
                    count is {count}
                </button>
                <p>
                    Edit <code>src/App.tsx</code> and save to test HMR
                </p>
            </div>
        </>
    )
}