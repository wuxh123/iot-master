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
      title: '设备模型',
      icon: 'cluster',
      open: false,
      children: [
        {
          title: '模型管理',
          router: 'model'
        },
        {
          title: '协议适配',
          router: 'tunnel'
        },
        {
          title: '变量管理',
          router: 'variable'
        },
        {
          title: '批量采集',
          router: 'batch'
        },
        {
          title: '定时任务',
          router: 'job'
        },
        {
          title: '自动策略',
          router: 'strategy'
        }
      ]
    },
    {
      title: '系统设置',
      icon: 'setting',
      open: false,
      children: [
        {
          title: '管理插件',
          router: 'plugin'
        },
        {
          title: '系统配置',
          router: 'setting'
        }
      ]
    }
  ];


  constructor() {

  }

  ngOnInit(): void {

  }
}
