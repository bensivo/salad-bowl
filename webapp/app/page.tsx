import Lobby from './lobby'
import styles from './page.module.css'

export default function Home() {
  return (
    <main className={styles.main}>
      <h1>Salad Bowl</h1>

      <Lobby></Lobby>
     
    </main>
  )
}
