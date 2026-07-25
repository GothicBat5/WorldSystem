import calendar

print("X To Exit.\n")

while True:
    yy = input("Enter year: ")
    if yy.lower() == 'x':
        print("\n**END**\n")
        break
    yy = int(yy)
    mm = input("Enter month: ")
    if mm.lower() == 'x':
        print("\n**END**\n")
        break
    mm = int(mm)
    print("\n")
    print(calendar.month(yy, mm))
    print("\n")
