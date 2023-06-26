# Salad Bowl Event Models

For documenting high level functionality, I use event modeling. You can read about it at https://eventmodeling.org/

The event models shown in this folder were created using Miro, and roughly represent the order of events as they occur in the game - as well as
which interfaces produce them, and which interfaces consume them.

Each event starts with a prefix, indicating its type. There are 4 types:
- trigger: Serialization of some real-life action
- request: Some request made by a user or service
- response: Response to a previously sent request event
- state: ECST-style state update. Carries the current state of an entity or collection of entities.
