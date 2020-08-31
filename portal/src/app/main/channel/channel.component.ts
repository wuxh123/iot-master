import {Component, OnInit, ViewContainerRef} from '@angular/core';
import {ApiService} from '../../api.service';
import {ChannelEditComponent} from '../channel-edit/channel-edit.component';
import {NzDrawerService, NzModalService} from 'ng-zorro-antd';

@Component({
  selector: 'app-channel',
  templateUrl: './channel.component.html',
  styleUrls: ['./channel.component.scss']
})
export class ChannelComponent implements OnInit {

  channels: [];

  constructor(private as: ApiService, private modal: NzModalService, private viewContainerRef: ViewContainerRef, private drawerService: NzDrawerService) {
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

  create(): void {
    this.drawerService.create({
      nzTitle: '创建',
      nzMaskClosable: false,
      nzWidth: 720,
      nzContent: ChannelEditComponent,
      nzContentParams: {
        channel: {}
      }
    });

  }

}
