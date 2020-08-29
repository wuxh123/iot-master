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
      title: '设置',
      icon: 'setting',
      open: false,
      children: [
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
