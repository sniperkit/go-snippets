#############################################################
#  Computer Project #10                                     #
#                                                           #
#   Python code that plays through a game of                #
#   Nine Men's Morris which is a game where you             #
#   need to complete three adjacent points to eliminate     #
#   an enemy piece and then continue on until the other     #
#   player has less than three pieces, this determines      #
#   the winner                                              #
#############################################################

import NMM #This is necessary for the project


BANNER = """
    __        _____ _   _ _   _ _____ ____  _ _ _ 
    \ \      / /_ _| \ | | \ | | ____|  _ \| | | |
     \ \ /\ / / | ||  \| |  \| |  _| | |_) | | | |
      \ V  V /  | || |\  | |\  | |___|  _ <|_|_|_|
       \_/\_/  |___|_| \_|_| \_|_____|_| \_(_|_|_)

"""


RULES = """
  _   _ _              __  __            _       __  __                 _     
 | \ | (_)_ __   ___  |  \/  | ___ _ __ ( )___  |  \/  | ___  _ __ _ __(_)___ 
 |  \| | | '_ \ / _ \ | |\/| |/ _ \ '_ \|// __| | |\/| |/ _ \| '__| '__| / __|
 | |\  | | | | |  __/ | |  | |  __/ | | | \__ \ | |  | | (_) | |  | |  | \__ \
 |_| \_|_|_| |_|\___| |_|  |_|\___|_| |_| |___/ |_|  |_|\___/|_|  |_|  |_|___/
                                                                                        
    The game is played on a grid where each intersection is a "point" and
    three points in a row is called a "mill". Each player has 9 pieces and
    in Phase 1 the players take turns placing their pieces on the board to 
    make mills. When a mill (or mills) is made one opponent's piece can be 
    removed from play. In Phase 2 play continues by moving pieces to 
    adjacent points. 
    
    The game is ends when a player (the loser) has less than three 
    pieces on the board.

"""


MENU = """

    Game commands (first character is a letter, second is a digit):
    
    xx        Place piece at point xx (only valid during Phase 1 of game)
    xx yy     Move piece from point xx to point yy (only valid during Phase 2)
    R         Restart the game
    H         Display this menu of commands
    Q         Quit the game
    
"""
        
def count_mills(board, player):
    """
        Counts how many mills are currently held by the player
        and returns the count. Function used before and after a
        placement is done
    """
    player_mills = 0
    for mill in board.MILLS:
        # Check if the first piece in the mill belongs to the player
        if all([board.points[m] == player for m in mill]):
            player_mills+=1
    return player_mills
            
def place_piece_and_remove_opponents(board, player, destination):
    """
        Used by a player if a mill is created, and 
        removes an opponent's piece. Takes three
        parameters, needs to raise a RuntimeError if 
        the placement is invalid. If placement is valid 
        it places the piece at destionation
    """
    mills_before = count_mills(board, player)
    # Check if the placement is valid
    if destination in board.points.keys():
        if board.points[destination] != " ":
            raise RuntimeError("Invalid command: Destination point already taken")
        else:
            board.assign_piece(player, destination)
            # Check to see if a new mill was created
            mills_after = count_mills(board, player)
            if mills_after > mills_before:
                print("A mill was formed!")
                print(board)
                remove_piece(board, player)
    else:
        raise RuntimeError("Invalid command: Not a valid point")
     
def move_piece(board, player, origin, destination):
    """
        This function moves a piece, takes in 4
        parameters. Raises a RuntimeError if movement 
        is invalid. If movement is valid it checks
        adjacency . Removes the player's piece from
        original point and calls functions to place
        the piece at new destination
    """
    mills_before = placed(board, player)
    if origin not in mills_before:
        raise RuntimeError("Invalid command: Origin point does not belong to player")
        #raise RuntimeError("Invalid command: Not a valid point")
    else:
        if destination not in board.ADJACENCY[origin]:
            #raise RuntimeError("Invalid command: Destination is not adjacent")
            raise RuntimeError("Invalid command: Not a valid point")
        else:
            if board.points[destination] == ' ':
                board.clear_place(origin)
                place_piece_and_remove_opponents(board, player, destination)
            else:
                #raise RuntimeError("Invalid command: Origin point does not belong to player")
                raise RuntimeError("Invalid command: Destination is not adjacent")



def points_not_in_mills(board, player):
    """
        Finds all points belonging to players that
        aren't in mills and returns them as
        an iterable data structure of either a list
        or set. use the list 'mills' in the Board class
    """
    player_mills = []
    for mill in board.MILLS:
        # Check if the first piece in the mill belongs to the player
        if all([board.points[m] == player for m in mill]):
            for m in mill:
                player_mills.append(m)
    not_mill_pieces = []
    for piece, piece_player in board.points.items():
        if player == piece_player:
            if piece not in player_mills:
                not_mill_pieces.append(piece)
    return not_mill_pieces

def placed(board,player):
    """
        Returns points of where the plyaer's pieces
        have been placed. Returns a set or a list
    """
    player_points = [k for k,v in board.points.items() if v==player]
    return player_points
    
def remove_piece(board, player):
    """
        This function removes a piece belonging to player
        on board. needs to determine if points are valid
        to remove (points_not_in_mills and placed def for help)
        loop until a valid piece is removed and handle the removal
        by calling clear_place method in Board class
    """
    player_pieces = placed(board,get_other_player(player))
    player_points_not_in_mills = points_not_in_mills(board,get_other_player(player))
    while True:
        try:
            piece_to_remove = input("Remove a piece at :> ").strip().lower()
            if piece_to_remove in player_pieces:
                # Only check if the piece isn't in a mill if there are no pieces that aren't
                if player_points_not_in_mills and piece_to_remove not in player_points_not_in_mills:
                    print("Invalid command: Point is in a mill")
                    print("Try again.")
                    raise RuntimeError("Invalid command: Point is in a mill")
                else:
                    # Remove the piece
                    board.clear_place(piece_to_remove)
                    return
            elif piece_to_remove not in board.points.keys():
                print("Invalid command: Not a valid point")
                print("Try again.")
                raise RuntimeError("Invalid command: Not a valid point")
            else:
                print("Invalid command: Point does not belong to player")
                print("Try again.")
                raise RuntimeError("Invalid command: Origin point does not belong to player")
        except RuntimeError as error_message:
            if error_message == "Invalid command: Point is in a mill" or error_message == "Invalid command: Origin point does not belong to player":
                pass
           
def is_winner(board, player):
    """
        Function used to decide if a game has been won
        a game is won if the other player has less than three
        pieces
    """
    new_player = get_other_player(player)
    opponent_pieces = 0
    for piece, piece_player in board.points.items():
        if (piece_player == new_player):
            opponent_pieces+=1
    return (opponent_pieces < 3)
   
def get_other_player(player):
    """
        Get's other player and returns it, use function
        any time you need to use the other player
    """
    return "X" if player == "O" else "O"
    
def main():
    #Loop so that we can start over on reset
    while True:
        #Setup stuff.
        print(RULES)
        print(MENU)
        board = NMM.Board()
        print(board)
        player = "X"
        placed_count = 0 # total of pieces placed by "X" or "O", includes pieces placed and then removed by opponent
        
        # PHASE 1
        print(player + "'s turn!")
        #placed = 0
        command = input("Place a piece at :> ").strip().lower()
        if command == "q":
            return
        print()
        if command == "h":
            print(MENU)
            print()
            command = input("Place a piece at :> ").strip().lower()
            if command == ' ':
                return 
        #Until someone quits or we place all 18 pieces...
        while command != 'q' and placed_count != 18:
            if command == 'r':
                main()
                return
            if command == "h":
                print(MENU)
                print()
                command = input("Place a piece at :> ").strip().lower()
                if command == "q":
                    return
            try:
                place_piece_and_remove_opponents(board, player, command)
                # Swap player
                player = get_other_player(player)
                placed_count+=1
            #Any RuntimeError you raise inside this try lands here
            except RuntimeError as error_message:
                print("{:s}\nTry again.".format(str(error_message)))
           
            #Prompt again
            print(board)
            print(player + "'s turn!")
            if placed_count < 18:
                command = input("Place a piece at :> ").strip().lower()
                if command == "q":
                    return
            else:
                print("**** Begin Phase 2: Move pieces by specifying two points")
                command = input("Move a piece (source,destination) :> ").strip().lower()
                if command == ' ':
                    return
            print()
        
        #Go back to top if reset
        if command == "q":
            return
        if command == 'r':
            main()
            return
        
        # PHASE 2 of game
        while command != 'q':
            if command == 'r':
                main()
                return
            # commands should have two points
            command = command.split()
            try:
                if len(command) > 1:
                    origin, destination = command
                    move_piece(board, player, origin, destination)
                    if is_winner(board, player):
                        print(BANNER)
                        return
                    else:
                        # Swap player
                        player = get_other_player(player)
                else:
                    raise RuntimeError("Invalid number of points")
            #Any RuntimeError you raise inside this try lands here
            except RuntimeError as error_message:
                print("{:s}\nTry again.".format(str(error_message)))         
            #Display and reprompt
            print(board)
            #display_board(board)
            print(player + "'s turn!")
            command = input("Move a piece (source,destination) :> ").strip().lower()
            if command == ' ':
                return

            print()
            
        #If we ever quit we need to return
        if command == 'q':
            return
    

            
if __name__ == "__main__":
    main()
