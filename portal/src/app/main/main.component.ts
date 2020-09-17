import {Component, OnInit} from '@angular/core';
import {MqttService} from '../mqtt.service';


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

  constructor(private mqtt: MqttService) {


  }

  ngOnInit(): void {
    this.mqtt.subscribe('/+/+/recv').subscribe(packet => {
      console.log('packet', packet);
    });
  }

}
