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
          router: 'connection'
        }
      ]
    },
    {
      title: '数据采集',
      icon: 'hdd',
      children: [
        {
          title: 'Modbus',
          router: 'modbus'
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
