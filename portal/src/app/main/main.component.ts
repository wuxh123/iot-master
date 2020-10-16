import {
  Component,
  OnInit,
} from '@angular/core';


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
      title: '设备管理',
      icon: 'api',
      children: [
        {
          title: '所有设备',
          router: 'device'
        },
        {
          title: '地图模式',
          router: 'device-map'
        },
        {
          title: '报警消息',
          router: 'device-alarm'
        },
        {
          title: '操作日志',
          router: 'device-log'
        },
      ]
    },
    {
      title: '数据通道',
      icon: 'api',
      children: [
        {
          title: '通道管理',
          router: 'tunnel'
        },
        {
          title: '连接管理',
          router: 'link'
        }
      ]
    },
    {
      title: '设备模型',
      icon: 'cluster',
      open: false,
      children: [
        {
          title: '模型管理',
          router: 'model'
        },
        {
          title: '元件管理',
          router: 'element'
        },
        {
          title: '协议管理',
          router: 'adapter'
        },
      ]
    },
    {
      title: '设置',
      icon: 'setting',
      open: false,
      children: [
        {
          title: '设置',
          router: 'setting'
        },
        {
          title: '数据备份',
          router: 'backup'
        },
        {
          title: '系统日志',
          router: 'logs'
        },
      ]
    }
  ];


  constructor() {

  }

  ngOnInit(): void {

  }
}
