<?php

use PHPMailer\PHPMailer\PHPMailer;
use PHPMailer\PHPMailer\SMTP;
use PHPMailer\PHPMailer\Exception;

/**
 * 盐
 * @param $key
 * @return string
 */
function salt($key)
{
    return sha1(pow($key, 2) * 3.14);
}

/**
 * 获取浏览器信息
 * @return string
 */
function getBrowser()
{
    $user_OSagent = $_SERVER['HTTP_USER_AGENT'];
    if (strpos($user_OSagent, "Maxthon") && strpos($user_OSagent, "MSIE")) {
        $visitor_browser = "Maxthon(Microsoft IE)";
    } elseif (strpos($user_OSagent, "Maxthon 2.0")) {
        $visitor_browser = "Maxthon 2.0";
    } elseif (strpos($user_OSagent, "Maxthon")) {
        $visitor_browser = "Maxthon";
    } elseif (strpos($user_OSagent, "Edge")) {
        $visitor_browser = "Edge";
    } elseif (strpos($user_OSagent, "Trident")) {
        $visitor_browser = "IE";
    } elseif (strpos($user_OSagent, "MSIE")) {
        $visitor_browser = "IE";
    } elseif (strpos($user_OSagent, "MSIE")) {
        $visitor_browser = "MSIE";
    } elseif (strpos($user_OSagent, "NetCaptor")) {
        $visitor_browser = "NetCaptor";
    } elseif (strpos($user_OSagent, "Netscape")) {
        $visitor_browser = "Netscape";
    } elseif (strpos($user_OSagent, "Chrome")) {
        $visitor_browser = "Chrome";
    } elseif (strpos($user_OSagent, "Lynx")) {
        $visitor_browser = "Lynx";
    } elseif (strpos($user_OSagent, "Opera")) {
        $visitor_browser = "Opera";
    } elseif (strpos($user_OSagent, "MicroMessenger")) {
        $visitor_browser = "WeiXinBrowser";
    } elseif (strpos($user_OSagent, "Konqueror")) {
        $visitor_browser = "Konqueror";
    } elseif (strpos($user_OSagent, "Mozilla/5.0")) {
        $visitor_browser = "Mozilla";
    } elseif (strpos($user_OSagent, "Firefox")) {
        $visitor_browser = "Firefox";
    } elseif (strpos($user_OSagent, "U")) {
        $visitor_browser = "Firefox";
    } else {
        $visitor_browser = "Other Browser";
    }
    return $visitor_browser;
}

/**
 * 获取操作系统信息
 * @return string
 */
function getOS()
{
    $OS = $_SERVER['HTTP_USER_AGENT'];
    if (preg_match('/win/i', $OS)) {
        $OS = 'Windows';
    } elseif (preg_match('/mac/i', $OS)) {
        $OS = 'MAC';
    } elseif (preg_match('/linux/i', $OS)) {
        $OS = 'Linux';
    } elseif (preg_match('/unix/i', $OS)) {
        $OS = 'Unix';
    } elseif (preg_match('/bsd/i', $OS)) {
        $OS = 'BSD';
    } else {
        $OS = 'Other';
    }
    return $OS;
}


/**
 * 发送邮件
 * @param $to
 * @param $content
 * @return bool
 * @throws \Exception
 */
function mailto($to, $content, $title)
{

    // Instantiation and passing `true` enables exceptions
    $mail = new PHPMailer(true);

    try {
        //Server settings
        $mail->SMTPDebug = SMTP::DEBUG_OFF;                      // Enable verbose debug output
        $mail->isSMTP();                                            // Send using SMTP
        $mail->Host = config('comment.host');                    // Set the SMTP server to send through
        $mail->SMTPAuth = true;                                   // Enable SMTP authentication
        $mail->Username = config('comment.username');                     // SMTP username
        $mail->Password = config('comment.password');                               // SMTP password
        $mail->SMTPSecure = PHPMailer::ENCRYPTION_SMTPS;         // Enable TLS encryption; `PHPMailer::ENCRYPTION_SMTPS` encouraged
        $mail->Port = 587;                                    // TCP port to connect to, use 465 for `PHPMailer::ENCRYPTION_SMTPS` above
        $mail->CharSet = 'utf-8';

        //Recipients
        $mail->setFrom('muxi_k_ing@163.com', 'Muxi_k');
        $mail->addAddress($to);     // Add a recipient

        // Content
        $mail->isHTML(true);                                  // Set email format to HTML
        $mail->Subject = $title;
        $mail->Body = $content;
        $mail->AltBody = $content;

        return $mail->send();
    } catch (Exception $e) {
        exception($mail->ErrorInfo, 1001);
    }
}

/**
 * 回复邮件模板
 * @param $nickname
 * @param $content
 * @param $reply
 * @param $repler
 * @param $link
 * @param $title
 * @return string
 */
function getReplyTemplate($nickname, $content, $reply, $repler, $link, $title)
{
    $template = <<<EOF
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <style>
        body {font-size: 18px;background: #cacaca;display: flex;align-items: center;justify-content: center; }
        .container {background-color: #fff;border-radius: 1em; margin: 10px}
        .title {border-radius: 1em 1em 0 0;padding: 1em 2em;color: #fff;background-image: linear-gradient(to right, rgb(3, 89, 146), #239dbd)}
        .box {padding: .5em 1.6em;max-width: 500px;}
        .content span a {color: #0d8ddb;text-decoration: none;}
        .content {border-radius: .4em}
        .content p{display: block;margin: 1em .2em;box-shadow: 0 0 .2em #ccc;padding: .3em;background-color: rgba(193, 184, 184, 0.05);}
    </style>
</head>
<body>
<div class="container">
    <div class="title">你在Muxik's Blog 的评论有了新的回复</div>
    <div class="box">
        <div class="content">
            <span class="content-title">{$nickname}, 你曾在文章 <a href="{$link}">《{$title}》</a>上发表评论：</span>
            <p>{$content}</p>
        </div>
        <div class="content">
            <span class="content-title">{$repler} 给你的回复如下:</span>
            <p>{$reply}</p>
        </div>
    </div>
</div>
</body>
</html>
EOF;

    return $template;
}

/**
 * 评论管理员模板
 * @param $nickname
 * @param $email
 * @param $link
 * @param $title
 * @param $content
 * @return string
 */
function getCommentTemplate($nickname, $email, $link, $title, $content)
{
    $template = <<<EOF
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <style>
        body {font-size: 18px;background: #cacaca;display: flex;align-items: center;justify-content: center; }
        .container {background-color: #fff;border-radius: 1em; margin: 10px}
        .title {border-radius: 1em 1em 0 0;padding: 1em 2em;color: #fff;background-image: linear-gradient(to right, rgb(3, 89, 146), #239dbd)}
        .box {padding: .5em 1.6em;max-width: 500px;}
        .content span a {color: #0d8ddb;text-decoration: none;}
        .content {border-radius: .4em}
        .content p{display: block;margin: 1em .2em;box-shadow: 0 0 .2em #ccc;padding: .3em;background-color: rgba(193, 184, 184, 0.05);}
    </style>
</head>
<body>
<div class="container">
    <div class="title">Muxik's Blog 有了新的评论</div>
    <div class="box">
        <div class="content">
            <span class="content-title">{$nickname}:{$email} 在文章 <a href="{$link}">{$title}</a>评论了:</span>
            <p>{$content}</p>
        </div>
    </div>
</div>
</body>
</html>
EOF;
    return $template;

}