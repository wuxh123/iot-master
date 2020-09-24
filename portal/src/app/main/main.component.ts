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
    router.events.subscribe((e: any) => {
      //console.log('router event', e)
      // if (this.container && this.container.createComponent
      //   && e.snapshot && e.snapshot.component && e.snapshot.children.length == 0) {
      //
      //   const factory = resolver.resolveComponentFactory(e.snapshot.component);
      //   this.container.clear();
      //   const ref = this.container.createComponent(factory);
      //   //TODO 第一次，container还未准备好
      //   //TODO 添加路由参数
      //   //TODO call ref.destroy() in ngOnDestroy
      //
      //
      // }
    });
    //TODO unsubscribe

  }

  tabs = [];

  onTabsChange(e): void {
    //TODO 实现标签页
    // 1、打开对应标签，如果未加载组件，则解析路由，导入组件，已加载则忽略
    // 2、路由切换，没有自动打开对应标签，可能情况：
    //  a、标签未创建或已经关闭，此时需要自动创建新标签，打开页面
    //  b、无效路由
    //

    //e.index

    // this.router.routerState.snapshot.root.firstChild.firstChild.firstChild...
    let route: any = this.router.routerState.snapshot.root;
    while (route.firstChild) {
      route = route.firstChild;
    }
    const factory = this.resolver.resolveComponentFactory(route.component);

    this.container.clear();
    const ref = this.container.createComponent(factory);

  }

  ngOnInit(): void {
    this.mqtt.subscribe('/+/+/recv').subscribe(packet => {
      console.log('packet', packet);
    });
  }

}
