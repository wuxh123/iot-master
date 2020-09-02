import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent implements OnInit {

  isCollapsed = false;

  menus = [
    {
      title: '控制台',
      icon: 'dashboard',
      open: true,
      children: [
        {
          title: '仪表盘',
          router: 'dash'
        },
      ]
    },
    {
      title: '数据通道',
      icon: 'api',
      children: [
        {
          title: '通道管理',
          router: 'channel'
        },
        {
          title: '连接管理',
          router: 'link'
        }
      ]
    },
    {
      title: '插件管理',
      icon: 'appstore',
      open: false,
      children: [
        {
          title: '已安装插件',
          router: 'plugin'
        },
        {
          title: '插件市场',
          router: 'plugin-store'
        }
      ]
    },
    {
      title: '系统设置',
      icon: 'setting',
      open: false,
      children: [
        {
          title: '用户管理',
          router: 'user'
        },
        {
          title: '修改密码',
          router: 'password'
        }
      ]
    }
  ];

  constructor() {
  }

  ngOnInit(): void {
  }

}
