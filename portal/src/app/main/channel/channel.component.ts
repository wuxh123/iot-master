import {Component, OnInit, ViewContainerRef} from '@angular/core';
import {ApiService} from '../../api.service';
import {ChannelEditComponent} from '../channel-edit/channel-edit.component';
import {NzDrawerService, NzModalService} from 'ng-zorro-antd';
import {ChannelDetailComponent} from "../channel-detail/channel-detail.component";

@Component({
  selector: 'app-channel',
  templateUrl: './channel.component.html',
  styleUrls: ['./channel.component.scss']
})
export class ChannelComponent implements OnInit {

  channels: [];

  constructor(private as: ApiService, private modal: NzModalService, private viewContainerRef: ViewContainerRef, private drawer: NzDrawerService) {
  }

  ngOnInit(): void {
    this.load();
  }

  load(): void {
    this.as.get('channels').subscribe(res => {
      if (res.ok) {
        this.channels = res.data;
      }
    });
  }

  edit(c?): void {
    this.drawer.create({
      nzTitle: c ? '编辑' : '创建',
      nzMaskClosable: false,
      nzWidth: 400,
      nzContent: ChannelEditComponent,
      nzContentParams: {
        channel: c || {}
      }
    });
  }

  detail(c): void {
    this.drawer.create({
      nzTitle: '详情',
      // nzWidth: 400,
      nzContent: ChannelDetailComponent,
      nzContentParams: {
        channel: c
      }
    });
  }

}
