import './app.less';
import { Homepage } from './pages/homepage';
import { useAppSelector } from './store/hooks';
import { routeSelectors } from './store/slices/route';

function App() {
  const route = useAppSelector(routeSelectors.route);

  const page = {
    '': <Homepage />,
  }[route] ?? <Homepage />

  return (
    <div className="app">
      {page}
    </div>
  )
}

export default App
