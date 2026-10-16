
<?php
session_start();
require "database.php";

$error = "";

if ($_SERVER["REQUEST_METHOD"] === "POST") 
{
    $username = trim($_POST["username"] ?? "");
    $password = $_POST["password"] ?? "";

    $stmt = $pdo->prepare("SELECT id, username, password FROM users WHERE username = ?");

    $stmt->execute([$username]);
    $user = $stmt->fetch(PDO::FETCH_ASSOC);

    if ($user && password_verify($password, $user["password"])) 
    {
        session_regenerate_id(true);

        $_SESSION["user_id"] = $user["id"];
        $_SESSION["username"] = $user["username"];

        header("Location: welcome.php");
        exit;
    }

    $error = "Invalid username or password.";
}
?>

<!DOCTYPE html>
<html lang="en">
<head>
    <title>Log In</title>
    <link rel="stylesheet" href="style.css">
</head>
<body class="bg">
    <main class="frame">
        <h1>Welcome back.</h1>
        <p class="subtitle">Sign in to continue.</p>

        <?php if (isset($_GET["registered"])): ?>
            <p class="success">Account created. You can log in now.</p>
        <?php endif; ?>

        <?php if ($error !== ""): ?>
            <p class="message"><?= htmlspecialchars($error) ?></p>
        <?php endif; ?>

        <form method="POST" action="index.php" class="form">
            <label for="username">Username</label>
            <input class="input" id="username"
                   type="text" name="username"
                   maxlength="50" required>

            <label for="password">Password</label>
            <input class="input" id="password"
                   type="password" name="password" required>

            <button class="button" type="submit">
                Sign in
            </button>
        </form>

        <p class="footer">
            Don't have an account?
            <a href="sign_up.php">Sign up</a>
        </p>
    </main>
</body>
</html>
