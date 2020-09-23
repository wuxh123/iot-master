import {Component, ComponentFactoryResolver, OnInit, ViewChild, ViewContainerRef} from '@angular/core';
import {MqttService} from '../mqtt.service';
import {Router} from "@angular/router";


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

  @ViewChild("container", {read: ViewContainerRef}) container: ViewContainerRef;

  constructor(private mqtt: MqttService, private router: Router, private resolver: ComponentFactoryResolver) {
    router.events.subscribe((e:any) => {
      console.log('router event', e)
      if (this.container && this.container.createComponent
        && e.snapshot && e.snapshot.component && e.snapshot.children.length==0) {

        const factory = resolver.resolveComponentFactory(e.snapshot.component);
        this.container.clear();
        const ref = this.container.createComponent(factory);
        //TODO 第一次，container还未准备好
        //TODO 添加路由参数
        //TODO call ref.destroy() in ngOnDestroy
      }
    });
    //TODO unsubscribe

  }

  ngOnInit(): void {
    this.mqtt.subscribe('/+/+/recv').subscribe(packet => {
      console.log('packet', packet);
    });
  }

}
