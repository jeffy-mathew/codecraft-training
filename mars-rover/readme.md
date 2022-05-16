TDD Exercise #1
The Mars Rover Kata (simplified)
Test-drive code to drive a rover over the idealised surface of Mars, represented as a 2-D grid of x, y coordinates (e.g., [4, 7]).
The rover can be “dropped” on to the grid at a specific set of coordinates, facing in one of four directions: North, East, South or West.
Sequences of instructions are sent to the rover telling it where to go. This sequence is a string of characters, each representing a single instruction:
 R = turn right
 L = turn left
 F = move on square forward
 B = move one square back
The rover ignores invalid instructions.
How to tackle this exercise:
1. Make a test list for your rover based on these requirements
2. Work through your test list on test case at a time, remembering to Red-Green-Refactor with
   each test
3. When refactoring, don’t forget that test code has to be maintained, too!