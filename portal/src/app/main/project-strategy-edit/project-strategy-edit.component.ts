import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzMessageService, NzModalRef} from 'ng-zorro-antd';

@Component({
  selector: 'app-project-strategy-edit',
  templateUrl: './project-strategy-edit.component.html',
  styleUrls: ['./project-strategy-edit.component.scss']
})
export class ProjectStrategyEditComponent implements OnInit {
  target = 'project/strategy';
  @Input() id = 0;

  data: any = {};

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
