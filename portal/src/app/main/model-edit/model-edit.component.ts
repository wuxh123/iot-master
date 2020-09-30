import {Component, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {ActivatedRoute, Router} from '@angular/router';

@Component({
  selector: 'app-model-edit',
  templateUrl: './model-edit.component.html',
  styleUrls: ['./model-edit.component.scss']
})
export class ModelEditComponent implements OnInit {
  title = '模型创建';
  id = 0;

  data: any = {};

  constructor(private as: ApiService, private router: Router, private routeInfo: ActivatedRoute) {
  }

  ngOnInit(): void {

    this.id = this.routeInfo.snapshot.params.id || 0;
    if (this.id > 0) {
      this.as.get('model/' + this.id).subscribe(res => {
        this.data = res.data;
      });
    }
  }

  submit(): void {
    if (this.data.id) {
      this.as.put('model/' + this.data.id, this.data).subscribe(res => {
        console.log(res);
        // TODO 修改成功
        this['closeTab']();
      });
    } else {
      this.as.post('model', this.data).subscribe(res => {
        console.log(res);
        // TODO 保存成功
        this['closeTab']();
      });
    }
  }
}
