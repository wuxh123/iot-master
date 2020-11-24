import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {ActivatedRoute} from '@angular/router';
import {TabRef} from '../tabs/tabs.component';
import {NzMessageService, NzModalRef} from 'ng-zorro-antd';

@Component({
  selector: 'app-project-edit',
  templateUrl: './project-edit.component.html',
  styleUrls: ['./project-edit.component.scss']
})
export class ProjectEditComponent implements OnInit {
  target = 'project';
  @Input() id = 0;

  data: any = {};

  content = "";

  constructor(private as: ApiService, private mr: NzModalRef, private ms: NzMessageService) {

  }

  ngOnInit(): void {
    if (this.id > 0) {
      this.as.get(this.target + '/' + this.id).subscribe(res => {
        this.data = res.data;
      });
    }
  }

  submit(): void {
    let uri = this.target;
    if (this.data.id) {
      uri += '/' + this.data.id;
    }
    this.as.post(uri, this.data).subscribe(res => {
      if (res.ok) {
        this.ms.success('保存成功');
        this.mr.close(res.data);
      }
    });
  }
}
