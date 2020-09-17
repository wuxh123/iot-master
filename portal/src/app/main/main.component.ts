import {Component, OnInit} from '@angular/core';

import * as mqtt from 'mqtt';

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

  constructor() {
    const client = mqtt.connect('ws://127.0.0.1:8080/api/mqtt');
    client.on('connect', data => {
      console.log('connect', data);
      client.subscribe('add', err => {
        console.log(err);
        if (!err) {
          client.publish('add', 'Hello mqtt');
        }
      });
      // client.subscribe('add');
      // client.publish('add', 'hello world');
    });
    client.on('message', (topic, message) => {
      console.log('message', topic, message);
    });
    client.on('close', () => {
      console.log('close');
    });
    client.on('offline', () => {
      console.log('offline');
    });
    client.on('disconnect', (data) => {
      console.log('disconnect');
    });
    client.on('error', (err) => {
      console.log('error', err);
    });

  }

  ngOnInit(): void {
  }

}
