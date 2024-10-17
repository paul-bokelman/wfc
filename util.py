class Point:
    def __init__(self, x, y):
        self.x: int = x
        self.y: int = y

    def out_of_grid(self, grid_size: int):
        if self.x >= grid_size or self.x < 0:
            return True
        if self.y >= grid_size or self.y < 0:
            return True

        return False
    
    def neighbor(self, direction: str):
        if direction == 'up':
            return Point(self.x, self.y + 1)
        if direction == 'down':
            return Point(self.x, self.y - 1)
        if direction == 'left':
            return Point(self.x - 1, self.y)
        if direction == 'right':
            return Point(self.x + 1, self.y)
        
        raise ValueError("Invalid direction")

    def __add__(self, other):
        return Point(self.x + other.x, self.y + other.y)
    
    def __sub__(self, other):
        return Point(self.x - other.x, self.y - other.y)
    
    def __eq__(self, other):
        return self.x == other.x and self.y == other.y

    def __str__(self) -> str:
        return f"({self.x}, {self.y})"