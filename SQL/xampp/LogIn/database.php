<?php

$host = "localhost";
$dbname = "login_system";
$dbuser = "root";
$dbpass = "";

try {

    $pdo = new PDO("mysql:host=$host; dbname=$dbname; charset=utf8mb4", $dbuser, $dbpass);
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
}
catch (PDOException $e) 
{
    error_log($e->getMessage());
    exit("Database connection failed.");
}
