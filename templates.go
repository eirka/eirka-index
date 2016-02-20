package main

// index template
const index = `[[define "index"]]<!doctype html>
<html ng-app="prim" ng-strict-di lang="en">
[[template "head" . ]]
<body>
<ng-include src="'pages/global.html'"></ng-include>
<div class="header">
[[template "header" . ]]
</div>
<div ng-view></div>
</body>
</html>[[end]]`

// head items
const head = `[[define "head"]]<head>
<base href="/[[ .base ]]">
<title data-ng-bind="page.title">[[ .title ]]</title>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<meta name="description" content="[[ .desc ]]" />
[[if .nsfw -]]<meta name="rating" content="adult" />
<meta name="rating" content="RTA-5042-1996-1400-1577-RTA" />[[- end]]
<link rel="stylesheet" href="/assets/prim/[[ .primcss ]]" />
<link rel="stylesheet" href="/assets/styles/[[ .style ]]" />
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.4.0/css/font-awesome.min.css">
<script src="/assets/prim/[[ .primjs ]]"></script>
[[template "angular" . ]][[template "headinclude" . ]]
</head>[[end]]`

// angular config
const angular = `[[define "angular"]]<script>angular.module('prim').constant('config',{
ib_id:[[ .ib ]],
title:'[[ .title ]]',
img_srv:'//[[ .imgsrv ]]',
api_srv:'//[[ .apisrv ]]',
csrf_token:'[[ .csrf ]]'
});
</script>[[end]]`

// site header
const header = `[[define "header"]]<div class="header_bar">
<div class="left">
<div class="nav_menu" ng-controller="NavMenuCtrl as navmenu">
<ul click-off="navmenu.close" ng-click="navmenu.toggle()" ng-mouseenter="navmenu.open()" ng-mouseleave="navmenu.close()">
<li class="n1"><a href><i class="fa fa-fw fa-bars"></i></a>
<ul ng-if="navmenu.visible">
[[template "navmenuinclude" . ]][[template "navmenu" . ]]
</ul>
</li>
</ul>
</div>
<div class="nav_items" ng-controller="NavItemsCtrl as navitems">
<ul>
<ng-include src="'pages/menus/nav.html'"></ng-include>
</ul>
</div>
</div>
<div class="right">
<div class="user_menu">
<div ng-if="!authState.isAuthenticated" class="login">
<a href="account" class="button-login">Sign in</a>
</div>
<div ng-if="authState.isAuthenticated" ng-controller="UserMenuCtrl as usermenu">
<ul click-off="usermenu.close" ng-click="usermenu.toggle()" ng-mouseenter="usermenu.open()" ng-mouseleave="usermenu.close()">
<li>
<div class="avatar avatar-medium">
<div class="avatar-inner">
<a href>
<img ng-src="{{authState.avatar}}" />
</a>
</div>
</div>
<ul ng-if="usermenu.visible">
<ng-include src="'pages/menus/user.html'"></ng-include>
</ul>
</li>
</ul>
</div>
</div>
<div class="site_logo">
<a href="/[[ .base ]]">
<img src="/assets/logo/[[ .logo ]]" title="[[ .title ]]" />
</a>
</div>
</div>
</div>[[end]]`

const navmenu = `[[define "navmenu"]][[ range $ib := .imageboards -]]
<li><a target="_self" href="//[[ $ib.Address ]]/">[[ $ib.Title ]]</a></li>
[[- end]][[end]]`
