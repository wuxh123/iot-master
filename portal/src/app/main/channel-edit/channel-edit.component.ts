import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {NzDrawerRef} from 'ng-zorro-antd';

@Component({
  selector: 'app-channel-edit',
  templateUrl: './channel-edit.component.html',
  styleUrls: ['./channel-edit.component.scss']
})
export class ChannelEditComponent implements OnInit {

  @Input() channel: any = {};

  constructor(private as: ApiService, private drawerRef: NzDrawerRef<string>) {
  }

  ngOnInit(): void {
    this.initChannel({});
    if (this.channel.id) {
      this.as.get('channel/' + this.channel.id).subscribe(res => {
        this.channel = res.data;
        this.initChannel(this.channel);
      });
    }
  }

  submit(): void {

    if (this.channel.id) {
      this.as.put('channel/' + this.channel.id, this.channel).subscribe(res => {
        console.log(res);
        // TODO 修改成功
        this.drawerRef.close(res.data);
      });
    } else {
      this.as.post('channel', this.channel).subscribe(res => {
        console.log(res);
        // TODO 保存成功
        this.drawerRef.close(res.data);
      });
    }
  }

  initChannel(item): void {
      item.net = item.net || { is_server: true, type: 'tcp', addr: ':1843'};
      item.register = item.register || {};
      item.heart_beat = item.heart_beat || {};

      this.channel = item;
  }
}
