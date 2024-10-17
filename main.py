from typing import TypedDict, cast
import random
from util import Point
from termcolor import colored, _types

class Rules(TypedDict):
    left: list[str]
    right: list[str]
    up: list[str]
    down: list[str]

class Cell(Point):
    def __init__(self, x: int, y: int, state: str | set[str]) -> None:
        super().__init__(x, y)
        self.state: str | set[str] = state
        self.collapsed = False # 

        if(isinstance(state, str)):
            self.collapsed = True

    def collapse(self):
        if self.collapsed:
            return

        self.state = random.choice(list(cast(set[str], self.state)))
        self.collapsed = True

    def set_possibilities(self, state: set[str]):
        self.state = state

        # collapse if singular option
        if len(state) == 1:
            self.collapse()

    def entropy(self):
        return len(self.state) if not self.collapsed else 0
            
class WaveFunctionCollapse:
    def __init__(self, tiles: dict[str, Rules], size: int = 5) -> None:
        self.tiles = tiles
        self.grid: list[list[Cell]] = [[Cell(x, y, set(tiles.keys())) for x in range(size)] for y in range(size)]
        self.size = size

    def _find_min_entropy(self):
        min_entropy_cell = None

        for row in self.grid:
            for cell in row:
                if not cell.collapsed:
                    entropy = cell.entropy()
                    if min_entropy_cell is None or entropy < min_entropy_cell.entropy():
                        min_entropy_cell = cell

        return min_entropy_cell
    
    # collapse and propagate until contradiction
    def collapse(self):
        while True:
            cell = self._find_min_entropy()

            if cell is None:  # Exit if no more collapsible cells
                break

            self.grid[cell.y][cell.x].collapse()
            self.propagate(self.grid[cell.y][cell.x])


    # propagates new possibilities
    def propagate(self, cell: Cell):
        # update surrounding possibilities
        for direction in ['up', 'down', 'left', 'right']:
            neighbor = cell.neighbor(direction)

            # update neighbor possibilities
            if not neighbor.out_of_grid(self.size) and not self.grid[neighbor.y][neighbor.x].collapsed:
                uncollapsed_neighbor = self.grid[neighbor.y][neighbor.x]
                possibilities = cast(set[str], uncollapsed_neighbor.state)
                possibilities.intersection_update(self.tiles[cast(str, cell.state)][direction]) 
                self.grid[neighbor.y][neighbor.x].set_possibilities(possibilities)

    def show(self):
        colors = {'A': 'red', 'B': 'blue', 'C': 'yellow'}
        for row in self.grid:
            for cell in row:

                if cell.collapsed:
                    print(colored('■', cast(_types.Color, colors[cast(str, cell.state)])), end=" ")
                else:
                    print("?", end=" ")
            print()

if __name__ == "__main__":

    tiles: dict[str, Rules] = {
        "A": {"left": ["A", "B"], "right": ["A", "B", "C"], "up": ["A", "B"], "down": ["A", "B"]},
        "B": {"left": ["A"], "right": ["A", "B"], "up": ["A", "B", "C"], "down": ["A", "B"]},
        "C": {"left": ["A", "B"], "right": ["A", "C"], "up": ["A", "C", "B"], "down": ["A"]},
    }

    wfc = WaveFunctionCollapse(tiles=tiles, size=10)
    wfc.collapse()
    wfc.show()
