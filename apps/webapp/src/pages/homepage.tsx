import axios from "axios";
import { useAppDispatch, useAppSelector } from "../store/hooks"
import { gamesActions, gamesSelectors } from "../store/slices/games/games-slice"

export function Homepage() {
    const games = useAppSelector(gamesSelectors.games);

    const dispatch = useAppDispatch();
    return (
        <>
            <h1>Salad bowl</h1>
            
            <button onClick={async () => {

                // TODO: turn the post + fetch action into its own thunk 'createGame'
                await new Promise((resolve) => setTimeout(resolve, Math.random() * 500));
                await axios.request({
                    method: 'POST',
                    url: 'http://localhost:8080/games'
                })

                await dispatch(gamesActions.fetchGames(null));
            }}>
                Create game
            </button>

            <table>
                <thead>
                    <tr>
                        <th>id</th>
                        <th></th>
                    </tr>
                </thead>

                <tbody>
                    {
                        games.map(game => (
                            <tr key={game.id}>
                                <td>{game.id}</td>
                                <td>
                                    <button>Join Game</button>
                                </td>
                            </tr>
                        ))
                    }
                </tbody>

            </table>

        </>
    )
}