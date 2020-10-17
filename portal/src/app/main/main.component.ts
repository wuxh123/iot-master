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
      title: '设备中心',
      icon: 'block',
      children: [
        {
          title: '设备管理',
          router: 'device'
        },
        {
          title: '地图模式',
          router: 'device-map'
        },
        {
          title: '操作日志',
          router: 'device-log'
        },
      ]
    },
    {
      title: '数据中心',
      icon: 'database',
      children: [
        {
          title: '历史记录',
          router: 'tunnel'
        },
        {
          title: '算法分析',
          router: 'link'
        },
      ]
    },
    {
      title: '报警中心',
      icon: 'alert',
      children: [
        {
          title: '报警记录',
          router: 'alert'
        },
        {
          title: '微信通知',
          router: 'alert-wechat'
        },
        {
          title: '邮件通知',
          router: 'alert-email'
        }
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
      title: '项目管理',
      icon: 'project',
      open: false,
      children: [
        {
          title: '项目管理',
          router: 'project'
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
      title: 'OTA升级',
      icon: 'cloud-upload',
      open: false,
      children: [
        {
          title: '固件管理',
          router: 'firmware'
        },
        {
          title: '升级日志',
          router: 'upgrade-log'
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
          title: '邮件发件箱',
          router: 'backup'
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
