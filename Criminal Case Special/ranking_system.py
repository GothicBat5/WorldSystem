
RANKS = [
    (1, "Officer"),
    (5, "Deputy"),
    (12, "Detective"),
    (20, "Corporal"),
    (32, "Sergeant"),
    (45, "Lieutenant"),
    (60, "Captain"),
    (75, "Major"),
    (100, "Inspector"),
    (125, "Lead Inspector"),
    (150, "Commander"),
    (175, "Deputy Chief"),
    (200, "Commissioner"),
    (250, "Sheriff"),
    (300, "Ranger"),
    (350, "Marshal"),
    (400, "Senior Trooper"),
    (500, "General"),
    (600, "Secret Agent"),
    (700, "Unit Chief"),
    (775, "Superintendent"),
    (850, "Vice Director")
]

ICONS = ["Bronze", "Silver", "Gold"]



#core functions

def get_rank_info(level):
    current_name = "Unknown"
    current_index = 0

    for i, (lvl, name) in enumerate(RANKS):
        if level >= lvl:
            current_name = name
            current_index = i
        else:
            break

    icon = ICONS[current_index % 3]
    return current_name, icon, current_index


def get_next_rank(level):
    for lvl, name in RANKS:
        if lvl > level:
            return lvl, name
    return None, None


def find_rank_by_name(rank_name):
    for lvl, name in RANKS:
        if name.lower() == rank_name.lower():
            return lvl, name
    return None, None


#da features

def check_level():
    level = int(input("Enter your level: "))
    name, icon, _ = get_rank_info(level)
    print(f"Level {level} → {name} (Icon: {icon})")


def show_by_icon():
    choice = input("\n (Bronze/Silver/Gold)\nChoose icon: ").capitalize()

    if choice not in ICONS:
        print("Invalid icon.")
        return

    print(f"\n--- {choice} Ranks ---")
    for i, (lvl, name) in enumerate(RANKS):
        icon = ICONS[i % 3]
        if icon == choice:
            print(f"Level {lvl} → {name}")


def show_progress():
    level = int(input("Enter your level: "))

    name, icon, _ = get_rank_info(level)
    next_lvl, next_name = get_next_rank(level)

    print(f"\nYou are: {name} ({icon})")

    if next_lvl:
        remaining = next_lvl - level
        print(f"Next Rank: {next_name}")
        print(f"Levels remaining: {remaining}")
    else:
        print("You reached the highest rank!")


def search_target_rank():
    level = int(input("Enter your current level: "))
    target = input("Enter target rank name: ")

    target_lvl, target_name = find_rank_by_name(target)

    if target_lvl is None:
        print("Rank not found.")
        return

    if level >= target_lvl:
        print(f"You already reached or passed {target_name}.")
    else:
        remaining = target_lvl - level
        print(f"You need {remaining} levels to reach {target_name}.")


#main

def main():
    while True:
        print("\n===** Rank System Menu **===")
        print("1. Check specific level")
        print("2. Show ranks by icon")
        print("3. Show progress")
        print("4. Search target rank")
        print("5. Exit")

        choice = input("Option: ")

        if choice == "1":
            check_level()

        elif choice == "2":
            show_by_icon()

        elif choice == "3":
            show_progress()

        elif choice == "4":
            search_target_rank()

        elif choice == "5":
            print("Goodbye!")
            break

        else:
            print("Invalid choice.")


#run program
if __name__ == "__main__":
    main()
