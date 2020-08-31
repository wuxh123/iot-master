import {Component, OnInit, ViewContainerRef} from '@angular/core';
import {ApiService} from '../../api.service';
import {ChannelEditComponent} from '../channel-edit/channel-edit.component';
import {NzModalService} from 'ng-zorro-antd';

@Component({
  selector: 'app-channel',
  templateUrl: './channel.component.html',
  styleUrls: ['./channel.component.scss']
})
export class ChannelComponent implements OnInit {

  channels: [];

  constructor(private as: ApiService, private modal: NzModalService, private viewContainerRef: ViewContainerRef) {
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
    this.modal.create({
      nzTitle: '创建',
      nzContent: ChannelEditComponent,
      nzViewContainerRef: this.viewContainerRef,
      nzGetContainer: () => document.body,
      nzComponentParams: {
        channel: {}
      }
    });

  }

}
