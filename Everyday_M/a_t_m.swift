import Foundation

class ATM {
    private var balance: Double
    
    init(initialBalance: Double) 
    {
        self.balance = initialBalance
    }
    
    func checkBalance() 
    {
        print("\nYour current balance is: \(balance)")
    }
    
    func deposit(amount: Double) 
    {
        guard amount > 0 else {
            print("\nDeposit amount must be greater than zero.")
            return
        }
        balance += amount
        print("\nYou deposited \(amount). New balance: \(balance)")
    }
    
    func withdraw(amount: Double) 
    {
        guard amount > 0 else {
            print("\nWithdrawal amount must be greater than zero.")
            return
        }
        if amount > balance 
        {
            print("\nInsufficient funds. Your balance is \(balance).")
        } 

        else {
            balance -= amount
            print("\nYou withdrew \(amount). New balance: \(balance)")
        }
    }
}

func showMenu() 
{
    print("""
    =========================
    Welcome to Swift ATM
    1. Check Balance
    2. Deposit
    3. Withdraw
    4. Exit
    =========================
    """)
}

let atm = ATM(initialBalance: 1000.0)

var running = true

while running {

    showMenu()
    print("Enter your choice: ", terminator: "")
    
    if let choice = readLine() 
    {
    
        switch choice {
        case "1":
            atm.checkBalance()
        case "2":
            print("Enter deposit amount: ", terminator: "")
            
            if let input = readLine(), let amount = Double(input) {
                atm.deposit(amount: amount)
            } 
            else {
                print("Invalid input.")
            }
        case "3":
            print("Enter withdrawal amount: ", terminator: "")
            
            if let input = readLine(), let amount = Double(input) {
                atm.withdraw(amount: amount)
            } 
            else {
                print("Invalid input.")
            }
        case "4":
            print("\nProgram Done!")
            running = false
        default:
            print("Invalid choice. Please try again.")
        }
    }
}
