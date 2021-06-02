class Quantity:
    
    TOO_HIGH = 2
    TOO_LOW = 0
    OKAY = 1

    def __init__(self, name="", min=0.0, max=0.0):
        self.name = name
        self.min = min
        self.max = max
        self.value = float()
        self.alert = ""
    
    def out_of_bounds(self):
        if self.value > self.max:
            return Quantity.TOO_HIGH
        elif self.value < self.min:
            return Quantity.TOO_LOW
        else:
            return Quantity.OKAY