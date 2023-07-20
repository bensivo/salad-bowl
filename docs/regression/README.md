# Regression Test

Regression test steps
1. Run unit test suite
2. Run e2e test suite
3. Run through manual tests


## Manual Tests
- Game Lobby
    - Create a game 
    - Join the game
        - Both windows should show both players
        - Each window should correctly identify itself as "me"
    - Join a team
        - Both windows should update
    - Refresh page, 
      - State should be persisted
    - Go back, 
      - Player should be shown disconnected
    - Go forward, 
      - Player should be shown connected
    - Start game, all windows update

- Word bank
    - Submit word
      - Updates UI
      - Updates others
    - Submit the same word
      - Shows error toast
    - Submit 3 words
      - Stops submissions
    - Refresh
      - State persists

- Etc.
  - Disconnect all players, wait 30 seconds, check localhost:8080/game deletes game