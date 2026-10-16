
<?php
session_start();
require "database.php";

$error = "";

if ($_SERVER["REQUEST_METHOD"] === "POST") 
{
    $username = trim($_POST["username"] ?? "");
    $password = $_POST["password"] ?? "";


    if ($username === "" || $password === "") 
    {
        $error = "Please fill in every field.";
    }

    elseif (strlen($password) < 8) 
    {
        $error = "Password must be at least 5 characters.";
    }

    else {
        $check = $pdo->prepare("SELECT id FROM users WHERE username = ?");
        $check->execute([$username]);

        if ($check->fetch()) 
        {
            $error = "That username is already taken.";
        
}        else {

            $hash = password_hash($password,PASSWORD_DEFAULT);

            $insert = $pdo->prepare("INSERT INTO users (username, password) VALUES (?, ?)");

            $insert->execute([$username, $hash]);

            header("Location: index.php?registered=1");
            exit;
        }
    }
}
?>

<!DOCTYPE html>
<html lang="en">
<head>
    <title>Sign Up</title>
    <link rel="stylesheet" href="style.css">
</head>
<body class="bg">
    <main class="frame">
        <h1>Create account</h1>
        <p class="subtitle">Get started with a new account.</p>

        <?php if ($error !== ""): ?>
            <p class="message"><?= htmlspecialchars($error) ?></p>
        <?php endif; ?>

        <form method="POST" action="sign_up.php" class="form">
            <label for="username">Username</label>
            <input class="input" id="username"
                   type="text" name="username"
                   maxlength="50" required>

            <label for="password">Password</label>
            <input class="input" id="password"
                   type="password" name="password"
                   minlength="5" required>

            <button class="button" type="submit">
                Sign up
            </button>
        </form>

        <p class="footer">
            Already registered?
            <a href="index.php">Log in</a>
        </p>
    </main>
</body>
</html>
