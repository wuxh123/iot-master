import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzDrawerRef} from 'ng-zorro-antd';

@Component({
  selector: 'app-model-edit',
  templateUrl: './model-edit.component.html',
  styleUrls: ['./model-edit.component.scss']
})
export class ModelEditComponent implements OnInit {

  @Input() model: any = {};


  constructor(private as: ApiService, private drawerRef: NzDrawerRef<string>) {
  }

  ngOnInit(): void {
    if (this.model.id) {
      this.as.get('model/' + this.model.id).subscribe(res => {
        this.model = res.data;
      });
    }
  }

  submit(): void {
    if (this.model.id) {
      this.as.put('model/' + this.model.id, this.model).subscribe(res => {
        console.log(res);
        // TODO 修改成功
        this.drawerRef.close(res.data);
      });
    } else {
      this.as.post('model', this.model).subscribe(res => {
        console.log(res);
        // TODO 保存成功
        this.drawerRef.close(res.data);
      });
    }
  }
}
