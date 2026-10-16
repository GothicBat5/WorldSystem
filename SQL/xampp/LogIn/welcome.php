<?php
session_start();

if (!isset($_SESSION["user_id"])) 
{
    header("Location: index.php");
    exit;
}

$username = htmlspecialchars($_SESSION["username"], ENT_QUOTES, "UTF-8");
?>

<!DOCTYPE html>
<html lang="en">
<head>
    <title>Welcome</title>
    <link rel="stylesheet" href="style.css">
</head>
<body class="bg">
    <main class="frame">
        <h1>Hello, <?= $username ?>!</h1>
        <p class="subtitle">Welcome back. You're logged in.</p>

        <a class="button button-link" href="log_out.php">
            Log out
        </a>
    </main>
</body>
</html>
