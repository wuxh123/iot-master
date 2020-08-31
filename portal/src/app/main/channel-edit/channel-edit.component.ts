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

  validateForm!: FormGroup;

  constructor(private as: ApiService, private fb: FormBuilder, private drawerRef: NzDrawerRef<string>) {
  }

  ngOnInit(): void {
    this.initForm({});
    if (this.channel.id) {
      this.as.get('channel/' + this.channel.id).subscribe(res => {
        this.channel = res.data;
        this.initForm(this.channel);
      });
    }

  }


  submit(): void {
    for (const i in this.validateForm.controls) {
      this.validateForm.controls[i].markAsDirty();
      this.validateForm.controls[i].updateValueAndValidity();
    }
    if (!this.validateForm.valid) {
      return;
    }

    if (this.channel.id) {
      this.as.put('channel/' + this.channel.id, this.validateForm.value).subscribe(res => {
        console.log(res);
        // TODO 修改成功
        this.drawerRef.close(res.data);
      });
    } else {
      this.as.post('channel', this.validateForm.value).subscribe(res => {
        console.log(res);
        // TODO 保存成功
        this.drawerRef.close(res.data);
      });
    }
  }

  initForm(item): void {
    if (!item.register)
      item.register = {};
    if (!item.heart_beat)
      item.heart_beat = {};

    this.validateForm = this.fb.group({
      name: [item.name, [Validators.required]],
      // tags: [item.tags],
      serial: [item.serial],
      net: [item.net, [Validators.required]],
      addr: [item.addr, [Validators.required]],
      is_server: [item.is_server],
      disabled: [item.disabled],
      'register.enable': [item.register.enable],
      'register.regex': [item.register.regex],
      'heart_beat.enable': [item.heart_beat.enable],
      'heart_beat.interval': [item.heart_beat.interval],
      'heart_beat.content': [item.heart_beat.content],
    });
  }
}
