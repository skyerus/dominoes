# Playing the game

```
docker pull skyerus/dominoes
docker run -p 8080:8080 skyerus/dominoes
```
Visit http://localhost:8080/

## Room for improvement

1. Garbage collection of stale/finished games
2. Ability for player to pick which side of the board to place their domino
3. A process that reveals each NPC's turn (instead of instantly placing their tiles)
